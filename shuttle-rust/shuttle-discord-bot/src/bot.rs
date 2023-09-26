use serenity::async_trait;
use serenity::model::channel::Message;
use serenity::model::gateway::Ready;
use serenity::prelude::*;
use tracing::info;

pub struct Bot {
    pub guild_id: String,
}

#[async_trait]
impl EventHandler for Bot {
    async fn message(&self, _ctx: Context, _new_message: Message) {
        todo!()
    }

    async fn ready(&self, _ctx: Context, _data_about_bot: Ready) {
        info!("The bot has successfully started and is now listening events...");
    }
}
