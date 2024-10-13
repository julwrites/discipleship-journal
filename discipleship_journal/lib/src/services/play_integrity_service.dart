import 'package:flutter/services.dart';

class PlayIntegrityService {
  static const platform = MethodChannel('io.tehj.developer.discipleship_journal/play_integrity');

  Future<bool> verifyDeviceIntegrity() async {
    try {
      final bool result = await platform.invokeMethod('verifyDeviceIntegrity');
      return result;
    } on PlatformException catch (e) {
      print('Error verifying device integrity: ${e.message}');
      return false;
    }
  }
}
