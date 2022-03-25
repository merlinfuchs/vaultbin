use sled::Error;

pub type DatabaseResult<T> = Result<T, DatabaseError>;

#[derive(Debug, thiserror::Error)]
pub enum DatabaseError {
    #[error("Database operation failed: {0}")]
    Sled(sled::Error),
    #[error("Decrypting database value failed: {0} ")]
    Encryption(aes_gcm::Error),
    #[error("Encoding database value failed: {0} ")]
    Encoding(rmp_serde::encode::Error),
    #[error("Decoding database value failed: {0} ")]
    Decoding(rmp_serde::decode::Error)
}

impl From<sled::Error> for DatabaseError {
    fn from(e: Error) -> Self {
        Self::Sled(e)
    }
}

impl From<aes_gcm::Error> for DatabaseError {
    fn from(e: aes_gcm::Error) -> Self {
        Self::Encryption(e)
    }
}

impl From<rmp_serde::encode::Error> for DatabaseError {
    fn from(e: rmp_serde::encode::Error) -> Self {
        Self::Encoding(e)
    }
}

impl From<rmp_serde::decode::Error> for DatabaseError {
    fn from(e: rmp_serde::decode::Error) -> Self {
        Self::Decoding(e)
    }
}