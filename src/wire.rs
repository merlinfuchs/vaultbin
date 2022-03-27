use actix_web::{HttpResponse, ResponseError};
use actix_web::body::BoxBody;
use actix_web::http::StatusCode;
use base64::DecodeError;
use log::error;
use serde::{Deserialize, Serialize};
use crate::database::DatabaseError;

pub type RouteResult<T> = Result<T, ApiError>;

#[derive(Debug, Serialize)]
pub struct ApiResponse<'a, T> {
    pub success: bool,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data: Option<T>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error: Option<&'a ApiError>,
}

impl<T: Serialize> ApiResponse<'_, T> {
    pub fn success(data: T) -> Self {
        Self {
            success: true,
            data: Some(data),
            error: None,
        }
    }
}

#[derive(Debug, Clone, thiserror::Error, Serialize)]
#[serde(tag = "code", rename_all = "snake_case")]
pub enum ApiError {
    #[error("Database operation failed")]
    DatabaseFailure,
    #[error("Unknown paste")]
    UnknownPaste,
    #[error("Invalid paste ID")]
    InvalidPasteId,
    #[error("Paste is too large")]
    PasteTooLarge,
    #[error("Failed to extract a ratelimit key")]
    RatelimitKeyExtractionFailed
}

impl ResponseError for ApiError {
    fn status_code(&self) -> StatusCode {
        match self {
            ApiError::DatabaseFailure => StatusCode::INTERNAL_SERVER_ERROR,
            ApiError::UnknownPaste => StatusCode::NOT_FOUND,
            ApiError::InvalidPasteId => StatusCode::NOT_FOUND,
            ApiError::PasteTooLarge => StatusCode::BAD_REQUEST,
            ApiError::RatelimitKeyExtractionFailed => StatusCode::INTERNAL_SERVER_ERROR
        }
    }

    fn error_response(&self) -> HttpResponse<BoxBody> {
        HttpResponse::build(self.status_code())
            .json(ApiResponse::<()> {
                success: false,
                data: None,
                error: Some(self),
            })
    }
}

impl From<DatabaseError> for ApiError {
    fn from(e: DatabaseError) -> Self {
        error!("{}", e);
        Self::DatabaseFailure
    }
}

impl From<DecodeError> for ApiError {
    fn from(_: DecodeError) -> Self {
        Self::InvalidPasteId
    }
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct CreatePasteRequestData {
    pub content: String,
    #[serde(default)]
    pub language: Option<String>,
    #[serde(default)]
    pub expiration: Option<u64>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct PasteResponseData {
    pub id: String,
    pub content: String,
    pub language: Option<String>,
    pub created_at: u64,
    pub view_count: u64
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct CreatePasteResponseData {
    #[serde(flatten)]
    pub paste: PasteResponseData,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct GetPasteResponseData {
    #[serde(flatten)]
    pub paste: PasteResponseData,
}