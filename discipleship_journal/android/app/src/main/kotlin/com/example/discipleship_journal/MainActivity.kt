package io.tehj.developer.discipleship_journal

import androidx.annotation.NonNull
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import com.google.android.play.core.integrity.IntegrityManagerFactory
import com.google.android.play.core.integrity.IntegrityTokenRequest
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch

class MainActivity: FlutterActivity() {
    private val CHANNEL = "io.tehj.developer.discipleship_journal/play_integrity"

    override fun configureFlutterEngine(@NonNull flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            if (call.method == "verifyDeviceIntegrity") {
                verifyDeviceIntegrity(result)
            } else {
                result.notImplemented()
            }
        }
    }

    private fun verifyDeviceIntegrity(result: MethodChannel.Result) {
        val integrityManager = IntegrityManagerFactory.create(this)
        val nonce = "YOUR_NONCE_HERE" // Replace with a unique nonce for each request

        val integrityTokenRequest = IntegrityTokenRequest.builder()
            .setNonce(nonce)
            .build()

        GlobalScope.launch(Dispatchers.Main) {
            try {
                val integrityTokenResponse = integrityManager.requestIntegrityToken(integrityTokenRequest)
                // Here you would typically send the token to your backend for verification
                // For this example, we'll just check if a response was received
                result.success(true)
            } catch (e: Exception) {
                result.error("INTEGRITY_CHECK_FAILED", e.message, null)
            }
        }
    }
}
