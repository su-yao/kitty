package io.legado.app.help.book

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import java.io.File

class DefaultTxtExportPathsTest {

    @Test
    fun `sanitized export defaults to download reading purify dir`() {
        val path = DefaultBookExportPaths.resolve(
            File("/storage/emulated/0/Download"),
            DefaultBookExportPaths.SANITIZED_DIR_NAME,
        )
        assertEquals("/storage/emulated/0/Download/阅读/净化", path)
        assertEquals("净化", File(path).name)
        assertEquals("阅读", File(path).parentFile?.name)
        assertEquals("Download", File(path).parentFile?.parentFile?.name)
    }

    @Test
    fun `original export defaults to download reading original dir`() {
        val path = DefaultBookExportPaths.resolve(
            File("/storage/emulated/0/Download"),
            DefaultBookExportPaths.ORIGINAL_DIR_NAME,
        )
        assertEquals("/storage/emulated/0/Download/阅读/原始", path)
        assertEquals("原始", File(path).name)
        assertEquals("阅读", File(path).parentFile?.name)
        assertEquals("Download", File(path).parentFile?.parentFile?.name)
    }

    @Test
    fun `legacy original txt cache path migrates to original dir`() {
        val migrated = DefaultBookExportPaths.migrateCachedPath(
            "/storage/emulated/0/Download/阅读/原TXT"
        )
        assertEquals("/storage/emulated/0/Download/阅读/原始", migrated)
        assertEquals(
            "/storage/emulated/0/Download/阅读/净化",
            DefaultBookExportPaths.migrateCachedPath("/storage/emulated/0/Download/阅读/净化")
        )
        assertEquals(null, DefaultBookExportPaths.migrateCachedPath(null))
        assertEquals("", DefaultBookExportPaths.migrateCachedPath(""))
        assertEquals("原TXT", DefaultBookExportPaths.migrateCachedPath("原TXT"))
    }

    @Test
    fun `epub export skips image audio video pdf and web file books`() {
        val source = readProjectFile("src/main/java/io/legado/app/help/book/BookExtensions.kt")
        assertTrue(source.contains("val Book.canExportEpub: Boolean"))
        assertTrue(source.contains("!isAudio && !isVideo && !isImage && !isPdf && !isWebFile"))
    }

    private fun readProjectFile(pathInApp: String): String =
        sequenceOf(File(pathInApp), File("app/$pathInApp"))
            .first { it.isFile }
            .readText()
            .replace("\r\n", "\n")
}
