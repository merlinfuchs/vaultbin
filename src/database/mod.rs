use std::time::{Duration, SystemTime, UNIX_EPOCH};
use base64::{CharacterSet, DecodeError};
use sha2::{Digest, Sha256};

use encryption::*;
pub use models::*;
pub use error::*;

use crate::DatabaseConfig;

mod models;
mod encryption;
mod error;

#[derive(Clone)]
pub struct Database {
    db: sled::Db,
}

pub fn encode_bytes_to_string(bytes: &[u8]) -> String {
    base64::encode_config(bytes, base64::Config::new(CharacterSet::UrlSafe, false))
}

pub fn decode_bytes_from_string(string: &str) -> Result<Vec<u8>, DecodeError> {
    base64::decode_config(string, base64::Config::new(CharacterSet::UrlSafe, false))
}

fn increment(old: Option<&[u8]>) -> Option<Vec<u8>> {
    let number = match old {
        Some(bytes) => {
            let array: [u8; 8] = bytes.try_into().unwrap();
            let number = u64::from_be_bytes(array);
            number + 1
        }
        None => 1,
    };

    Some(number.to_be_bytes().to_vec())
}

impl Database {
    pub fn new(config: &DatabaseConfig) -> DatabaseResult<Self> {
        let db = sled::Config::new()
            .path(&config.path)
            .cache_capacity(config.cache_size)
            .open()?;

        Ok(Self { db })
    }

    fn hash_bytes(bytes: &[u8]) -> Vec<u8> {
        let mut hasher = Sha256::new();
        hasher.update(&bytes);
        hasher.finalize().to_vec()
    }

    pub fn get_paste(&self, paste_id: &[u8]) -> DatabaseResult<Option<PasteModel>> {
        let hashed_id = Self::hash_bytes(paste_id);

        let res = match self.db.get(hashed_id.as_slice())? {
            Some(raw) => {
                let model = rmp_serde::from_slice(&raw)?;
                Some(decrypt_paste_model(paste_id, model)?)
            }
            None => None
        };
        Ok(res)
    }

    pub fn insert_paste(&self, model: PasteModel) -> DatabaseResult<[u8; 16]> {
        let key_bytes = generate_key();

        let hashed_id = Self::hash_bytes(&key_bytes);

        let encrypted = encrypt_paste_model(&key_bytes, model)?;
        let value = rmp_serde::to_vec(&encrypted)?;
        self.db.insert(&hashed_id, value)?;
        Ok(key_bytes)
    }

    pub fn count_paste_view(&self, paste_id: &[u8]) -> DatabaseResult<u64> {
        let tree = self.db.open_tree("views")?;

        let hashed_id = Self::hash_bytes(paste_id);

        let res = tree.update_and_fetch(hashed_id, increment)?.unwrap();
        let array: [u8; 8] = res.to_vec().try_into().unwrap();
        Ok(u64::from_be_bytes(array))
    }

    pub fn set_paste_expiration(&self, paste_id: &[u8], expiration: Duration) -> DatabaseResult<()> {
        let tree = self.db.open_tree("expirations")?;
        let hashed_id = Self::hash_bytes(paste_id);

        let expires_at = SystemTime::now().duration_since(UNIX_EPOCH).unwrap() + expiration;
        // this could theoretically result in collisions but it's very unlikely
        let expiry_key = expires_at.as_nanos().to_be_bytes();

        tree.insert(expiry_key, hashed_id)?;
        Ok(())
    }

    pub fn enforce_paste_expirations(&self) -> DatabaseResult<()> {
        let tree = self.db.open_tree("expirations")?;

        let start_expiry = 0u128.to_be_bytes();
        let unix_now = SystemTime::now().duration_since(UNIX_EPOCH).unwrap();
        let end_expiry = unix_now.as_nanos().to_be_bytes();

        for (_, paste_id) in tree.range(start_expiry..end_expiry).flatten() {
            self.db.remove(paste_id)?;
        }

        Ok(())
    }
}