use std::time::Duration;
use actix_web::{App, HttpServer, web, middleware};

use crate::config::*;
use crate::routes::*;

mod config;
mod routes;
mod database;
mod wire;

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
