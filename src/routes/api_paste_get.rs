use actix_web::{get, HttpResponse, Responder, web};

use crate::database::{Database, decode_bytes_from_string};
use crate::wire::{ApiError, ApiResponse, GetPasteResponseData, PasteResponseData, RouteResult};

#[get("/api/pastes/{paste_id}")]
pub async fn route_api_paste_get(db: web::Data<Database>, paste_id: web::Path<String>) -> RouteResult<impl Responder> {
    let paste_id = paste_id.into_inner();
    let decoded_paste_id = decode_bytes_from_string(&paste_id)?;

    let model = db.get_paste(&decoded_paste_id)?
        .ok_or(ApiError::UnknownPaste)?;

    let view_count = db.count_paste_view(&decoded_paste_id)?;

    let paste = PasteResponseData {
        id: paste_id,
        content: model.content,
        language: model.language,
        created_at: model.created_at,
        view_count,
    };
    Ok(HttpResponse::Ok()
        .json(ApiResponse::success(GetPasteResponseData { paste })))
}
