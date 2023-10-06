use axum::Router;
use axum::routing::get;
use shuttle_secrets::SecretStore;

async fn hello_world() -> &'static str {
    "Hello, world!"
}

#[shuttle_runtime::main]
async fn axum(
    #[shuttle_secrets::Secrets] secrets: SecretStore,
) -> shuttle_axum::ShuttleAxum {

    // Getting secrets from our SecretsStore - safe to unwrap as they're required for the app to work
    let oauth_id = secrets.get("OAUTH_CLIENT_ID").unwrap();
    let oauth_secret = secrets.get("OAUTH_CLIENT_SECRET").unwrap();

    let router = Router::new().route("/", get(hello_world));

    Ok(router.into())
}