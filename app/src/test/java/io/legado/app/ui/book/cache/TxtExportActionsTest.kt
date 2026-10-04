package io.legado.app.ui.book.cache

import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test
import java.io.File

class TxtExportActionsTest {

    @Test
    fun `cache export menu starts purified and original txt epub actions`() {
        val activity = readProjectFile(CACHE_ACTIVITY)
        val adapter = readProjectFile(CACHE_ADAPTER)
        val menu = readProjectFile(CACHE_MENU)
        val service = readProjectFile(EXPORT_SERVICE)

        assertTrue(adapter.contains("fun exportOptions(position: Int, anchor: View)"))
        assertTrue(adapter.contains("callBack.exportOptions(holder.layoutPosition, it)"))
        assertFalse(adapter.contains("fun export(position: Int)"))

        assertTrue(activity.contains("item(getString(R.string.export_sanitized_txt), \"sanitized\")"))
        assertTrue(activity.contains("item(getString(R.string.download_original_txt), \"original\")"))
        assertTrue(activity.contains("\"sanitized\" -> exportBooks(position, sanitized = true)"))
        assertTrue(activity.contains("\"original\" -> exportBooks(position, sanitized = false)"))
        assertTrue(activity.contains("R.id.menu_export_all_sanitized_txt -> exportAllBooks(sanitized = true)"))
        assertTrue(activity.contains("R.id.menu_download_all_original_txt -> exportAllBooks(sanitized = false)"))
        assertTrue(activity.contains("DefaultBookExportPaths.sanitized()"))
        assertTrue(activity.contains("DefaultBookExportPaths.original()"))
        assertTrue(activity.contains("DefaultBookExportPaths.migrateCachedPath"))
        assertTrue(activity.contains("putExtra(\"exportUseReplace\", sanitized)"))
        assertTrue(activity.contains("putExtra(\"exportType\", \"txt_epub\")"))

        assertTrue(menu.contains("android:id=\"@+id/menu_export_all_sanitized_txt\""))
        assertTrue(menu.contains("android:id=\"@+id/menu_download_all_original_txt\""))

        assertTrue(service.contains("val useReplace: Boolean? = null"))
        assertTrue(service.contains("\"txt\" -> exportTxt(exportConfig.path, book, exportConfig.useReplace)"))
        assertTrue(service.contains("\"txt_epub\" -> {"))
        assertTrue(service.contains("if (book.canExportEpub)"))
        assertTrue(service.contains("runCatching {"))
        assertTrue(service.contains("EPUB失败，已跳过"))
        assertTrue(service.contains("exportEpub(\n                                        exportConfig.path,\n                                        book,\n                                        exportConfig.useReplace"))
        assertTrue(service.contains("intent.hasExtra(\"exportUseReplace\")"))
        assertTrue(service.contains("private suspend fun getAllContents(\n        book: Book,\n        useReplace: Boolean,"))
        assertTrue(service.contains("private suspend fun exportEpub(path: String, book: Book, useReplace: Boolean? = null)"))
        assertTrue(service.contains("setEpubContent(contentModel, book, epubBook, useReplace)"))
    }

    private fun readProjectFile(pathInApp: String): String =
        sequenceOf(File(pathInApp), File("app/$pathInApp"))
            .firstOrNull(File::isFile)
            ?.readText()
            ?.replace("\r\n", "\n")
            .orEmpty()

    private companion object {
        const val CACHE_ACTIVITY = "src/main/java/io/legado/app/ui/book/cache/CacheActivity.kt"
        const val CACHE_ADAPTER = "src/main/java/io/legado/app/ui/book/cache/CacheAdapter.kt"
        const val CACHE_MENU = "src/main/res/menu/book_cache.xml"
        const val EXPORT_SERVICE = "src/main/java/io/legado/app/service/ExportBookService.kt"
    }
}
