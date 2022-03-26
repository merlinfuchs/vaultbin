use actix_web::{get, HttpResponse, Responder, web};

use crate::database::{Database, decode_bytes_from_string};
use crate::wire::{ApiError, RouteResult};

#[get("/api/pastes/{paste_id}/raw")]
pub async fn route_api_paste_get_raw(db: web::Data<Database>, paste_id: web::Path<String>) -> RouteResult<impl Responder> {
    let paste_id = paste_id.into_inner();
    let decoded_paste_id = decode_bytes_from_string(&paste_id)?;
    let model = db.get_paste(&decoded_paste_id)?
        .ok_or(ApiError::UnknownPaste)?;

    db.count_paste_view(&decoded_paste_id)?;
    Ok(HttpResponse::Ok()
        .append_header(("Content-Type", "text/plain"))
        .body(model.content))
}
