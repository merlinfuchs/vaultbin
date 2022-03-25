use actix_web::{get, HttpResponse, Responder, web};
use base64::CharacterSet;

use crate::database::Database;
use crate::wire::{ApiError, ApiResponse, GetPasteResponseData, PasteResponseData, RouteResult};

#[get("/api/pastes/{paste_id}")]
pub async fn route_api_paste_get(db: web::Data<Database>, paste_id: web::Path<String>) -> RouteResult<impl Responder> {
    let paste_id = paste_id.into_inner();
    let decode_paste_id = base64::decode_config(&paste_id, base64::Config::new(CharacterSet::UrlSafe, false))?;

    let model = db.get_paste(&decode_paste_id)?
        .ok_or(ApiError::UnknownPaste)?;

    let view_count = db.count_paste_view(&decode_paste_id)?;

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
