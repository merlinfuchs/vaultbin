use serde::{Deserialize, Serialize};

fn _is_default<T: Default + PartialEq>(t: &T) -> bool {
    *t == Default::default()
}

#[derive(Clone, Debug)]
pub struct PasteModel {
    pub language: Option<String>,
    pub created_at: u64,
    pub content: String
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct EncryptedPasteModel {
    pub language: Option<String>,
    pub created_at: u64,
    pub content: Vec<u8>,
    pub nonce: [u8; 12],
}