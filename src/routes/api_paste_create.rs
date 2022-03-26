use std::time::{Duration, SystemTime, UNIX_EPOCH};

use actix_web::{HttpResponse, post, Responder, web};

use crate::database::{Database, encode_bytes_to_string, PasteModel};
use crate::VaultbinConfig;
use crate::wire::{ApiError, ApiResponse, CreatePasteRequestData, CreatePasteResponseData, PasteResponseData, RouteResult};

#[post("/api/pastes")]
pub async fn route_api_paste_create(db: web::Data<Database>, data: web::Json<CreatePasteRequestData>, config: web::Data<VaultbinConfig>) -> RouteResult<impl Responder> {
    let data = data.into_inner();

    if data.content.len() > config.max_paste_size {
        return Err(ApiError::PasteTooLarge);
    }

    let model = PasteModel {
        content: data.content,
        language: data.language,
        created_at: SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs(),
    };

    let id = db.insert_paste(model.clone())?;
    let expiry = data.expiration.unwrap_or(config.max_expiration).min(config.max_expiration);
    db.set_paste_expiration(&id, Duration::from_secs(expiry))?;

    let encoded_id = encode_bytes_to_string(&id);
    let paste = PasteResponseData {
        id: encoded_id,
        content: model.content,
        language: model.language,
        created_at: model.created_at,
        view_count: 0,
    };
    Ok(HttpResponse::Ok()
        .json(ApiResponse::success(CreatePasteResponseData { paste })))
}
