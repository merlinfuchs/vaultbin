use std::net::IpAddr;
use std::time::Duration;

use actix_governor::{Governor, GovernorConfigBuilder, KeyExtractor};
use actix_web::{App, HttpServer, middleware, web};
use actix_web::dev::ServiceRequest;
use actix_web::http::Method;

use crate::config::*;
use crate::routes::*;
use crate::wire::ApiError;

mod config;
mod routes;
mod database;
mod wire;

#[derive(Clone, Copy)]
struct RatelimitKeyExtractor {
    pub reverse_proxy: bool,
}

impl KeyExtractor for RatelimitKeyExtractor {
    type Key = IpAddr;
    type KeyExtractionError = ApiError;

    fn extract(&self, req: &ServiceRequest) -> Result<Self::Key, Self::KeyExtractionError> {
        let res = match self.reverse_proxy {
            true => req
                .connection_info()
                .realip_remote_addr()
                .ok_or(ApiError::RatelimitKeyExtractionFailed)?
                .parse()
                .map_err(|_| ApiError::RatelimitKeyExtractionFailed)?,
            false => req.peer_addr()
                .ok_or(ApiError::RatelimitKeyExtractionFailed)?
                .ip()
        };
        Ok(res)
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let config = VaultbinConfig::new()
        .expect("parsing config");

    let database = database::Database::new(&config.database)
        .expect("opening database");

    let task_db = database.clone();
    actix_web::rt::spawn(async move {
        loop {
            let _ = task_db.enforce_paste_expirations();
            actix_web::rt::time::sleep(Duration::from_secs(60)).await;
        }
    });

    let governor_conf = GovernorConfigBuilder::default()
        .per_second(config.ratelimit.per_second)
        .burst_size(config.ratelimit.burst_size)
        .methods(vec![Method::POST])
        .key_extractor(RatelimitKeyExtractor {
            reverse_proxy: config.ratelimit.reverse_proxy
        })
        .finish()
        .unwrap();

    let port = config.port;
    let host = config.host.clone();

    HttpServer::new(move || {
        let app = App::new();

        #[cfg(feature = "cors")]
            let app = app.wrap(actix_cors::Cors::default()
            .allowed_methods(["GET", "POST", "OPTIONS"])
            .allow_any_header()
            .allow_any_origin());

        app
            .wrap(middleware::Compress::default())
            .wrap(Governor::new(&governor_conf))
            .app_data(web::Data::new(config.clone()))
            .app_data(web::Data::new(database.clone()))
            .service(route_api_paste_create)
            .service(route_api_paste_get)
            .service(route_api_paste_get_raw)
            .service(route_api_stats)
            .service(route_serve_frontend)
    })
        .bind((host.as_str(), port))?
        .run()
        .await
}
