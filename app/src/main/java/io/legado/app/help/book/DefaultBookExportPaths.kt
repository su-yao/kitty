package io.legado.app.help.book

import android.os.Environment
import java.io.File

object DefaultBookExportPaths {

    const val APP_DIR_NAME = "阅读"
    const val SANITIZED_DIR_NAME = "净化"
    const val ORIGINAL_DIR_NAME = "原始"

    fun sanitized(): String = path(SANITIZED_DIR_NAME)

    fun original(): String = path(ORIGINAL_DIR_NAME)

    fun resolve(downloadsDir: File, dirName: String): String =
        File(File(downloadsDir, APP_DIR_NAME), dirName).absolutePath

    fun migrateCachedPath(cached: String?): String? {
        if (cached.isNullOrEmpty()) return cached
        val file = File(cached)
        if (file.name != LEGACY_ORIGINAL_DIR_NAME) return cached
        val parent = file.parentFile ?: return cached
        return File(parent, ORIGINAL_DIR_NAME).absolutePath
    }

    private fun path(dirName: String): String {
        val downloads = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_DOWNLOADS)
        return resolve(downloads, dirName)
    }

    private const val LEGACY_ORIGINAL_DIR_NAME = "原TXT"
}
