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

        assertTrue(adapter.contains("fun export(position: Int)"))
        assertTrue(adapter.contains("fun exportSanitized(position: Int)"))
        assertTrue(adapter.contains("fun exportOriginal(position: Int)"))
        assertTrue(adapter.contains("callBack.exportSanitized(holder.layoutPosition)"))
        assertTrue(adapter.contains("callBack.exportOriginal(holder.layoutPosition)"))
        assertTrue(adapter.contains("callBack.export(holder.layoutPosition)"))
        assertFalse(adapter.contains("fun exportOptions(position: Int, anchor: View)"))

        assertTrue(activity.contains("override fun exportSanitized(position: Int)"))
        assertTrue(activity.contains("override fun exportOriginal(position: Int)"))
        assertTrue(activity.contains("exportBooks(position, sanitized = true)"))
        assertTrue(activity.contains("exportBooks(position, sanitized = false)"))
        assertTrue(activity.contains("R.id.menu_export_all_sanitized_txt -> exportAllBooks(sanitized = true)"))
        assertTrue(activity.contains("R.id.menu_download_all_original_txt -> exportAllBooks(sanitized = false)"))
        assertTrue(activity.contains("DefaultBookExportPaths.sanitized()"))
        assertTrue(activity.contains("DefaultBookExportPaths.original()"))
        assertTrue(activity.contains("DefaultBookExportPaths.migrateCachedPath"))
        assertTrue(activity.contains("FileUtils.createFolderIfNotExist(cached)"))
        assertTrue(activity.contains("putExtra(\"exportUseReplace\", sanitized)"))
        assertTrue(activity.contains("putExtra(\"exportType\", \"txt_epub\")"))
        assertTrue(adapter.contains("tvExportSanitized.setOnClickListener"))
        assertTrue(adapter.contains("tvExportOriginal.setOnClickListener"))
        assertTrue(adapter.contains("tvExport.setOnClickListener"))
        assertFalse(activity.contains("if (book.isLocal) return"))

        val itemLayout = readProjectFile(CACHE_ITEM)
        assertTrue(itemLayout.contains("android:id=\"@+id/tv_export_sanitized\""))
        assertTrue(itemLayout.contains("android:id=\"@+id/tv_export_original\""))
        assertTrue(itemLayout.contains("android:id=\"@+id/tv_export\""))
        assertTrue(itemLayout.contains("@string/export_sanitized_short"))
        assertTrue(itemLayout.contains("@string/download_original_short"))

        assertTrue(menu.contains("android:id=\"@+id/menu_export_all_sanitized_txt\""))
        assertTrue(menu.contains("android:id=\"@+id/menu_download_all_original_txt\""))

        val bookshelfMenu = readProjectFile(BOOKSHELF_MENU)
        assertTrue(bookshelfMenu.contains("android:id=\"@+id/menu_download\""))
        assertTrue(bookshelfMenu.contains("@string/cache_export"))
        assertTrue(
            Regex(
                """android:id="@\+id/menu_download"[\s\S]*?app:showAsAction="always""""
            ).containsMatchIn(bookshelfMenu)
        )

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
        const val CACHE_ITEM = "src/main/res/layout/item_download.xml"
        const val BOOKSHELF_MENU = "src/main/res/menu/main_bookshelf.xml"
        const val EXPORT_SERVICE = "src/main/java/io/legado/app/service/ExportBookService.kt"
    }
}
