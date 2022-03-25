use actix_web::{get, HttpResponse, Responder, web};
use base64::CharacterSet;

use crate::wire::{ApiError, RouteResult};

use crate::database::Database;

#[get("/api/pastes/{paste_id}/raw")]
pub async fn route_api_paste_get_raw(db: web::Data<Database>, paste_id: web::Path<String>) -> RouteResult<impl Responder> {
    let paste_id = paste_id.into_inner();
    let decode_paste_id = base64::decode_config(&paste_id, base64::Config::new(CharacterSet::UrlSafe, false))?;
    let model = db.get_paste(&decode_paste_id)?
        .ok_or(ApiError::UnknownPaste)?;

    db.count_paste_view(&decode_paste_id)?;
    Ok(HttpResponse::Ok()
        .append_header(("Content-Type", "text/plain"))
        .body(model.content))
}
