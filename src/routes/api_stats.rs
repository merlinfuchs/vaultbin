use actix_web::{get, Responder, web};

use crate::wire::RouteResult;
use crate::VaultbinConfig;

#[get("/api/stats")]
pub async fn route_api_stats(_config: web::Data<VaultbinConfig>) -> RouteResult<impl Responder> {
    Ok("")
}
