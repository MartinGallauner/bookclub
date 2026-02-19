import 'package:flutter/material.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

class BarcodeScannerSheet extends StatefulWidget {
  const BarcodeScannerSheet({super.key});

  @override
  State<StatefulWidget> createState() => _BarcodeScannerSheetState();
}

class _BarcodeScannerSheetState extends State<BarcodeScannerSheet> {
  @override
  Widget build(BuildContext context) {
    bool _hasScanned = false;
    return Stack(
      children: [
        Positioned.fill(child: MobileScanner()),
        Positioned(
          top: 0,
          left: 0,
          child: SafeArea(
            child: CircleAvatar(
              backgroundColor: Colors.black,
              radius: 20,
              child: IconButton(
                icon: Icon(Icons.close),
                onPressed: () {
                  Navigator.of(context).pop(); // close without result
                },
              ),
            ),
          ),
        ),
      ],
    );
  }
}
