package me.clement_casse.jetbrains_webapp_playground

import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.server.testing.*
import org.jetbrains.exposed.sql.Database
import org.junit.Test
import kotlin.test.assertEquals

class ApplicationTest {
    private val testDb = Database.connect("jdbc:h2:mem:test;DB_CLOSE_DELAY=-1", driver = "org.h2.Driver")

    @Test
    fun helloTest() = testApplication {
        application {
            mainWithDependencies(testDb)
        }
        val helloResp = client.get("/")
        assertEquals("Hello World!", helloResp.bodyAsText())
    }
}