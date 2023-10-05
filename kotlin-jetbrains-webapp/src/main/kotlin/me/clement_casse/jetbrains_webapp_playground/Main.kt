package me.clement_casse.jetbrains_webapp_playground

import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.callloging.*
import io.ktor.server.plugins.conditionalheaders.*
import io.ktor.server.plugins.defaultheaders.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import org.jetbrains.exposed.sql.Database
import java.io.File

fun main() {
    embeddedServer(Netty, port = 8080, module = Application::main).start(wait = true)
}

fun Application.main() {
    val h2DbFile = File("build/db")
    val db = Database.connect(
        "jdbc:h2:file:${h2DbFile.canonicalFile.absolutePath}",
        driver = "org.h2.Driver"
    )
    mainWithDependencies(db)
}

fun Application.mainWithDependencies(db: Database){
    install(DefaultHeaders)      // Set `Server` & `Date` headers into each response
    install(ConditionalHeaders)  // Avoids sending the body of content if it has not changed since the last request.
    install(CallLogging)         // Log each requests

    // Defining the HTTP Router of the application
    routing {
        get("/") {
            call.respond("Hello World!")
        }
    }
}