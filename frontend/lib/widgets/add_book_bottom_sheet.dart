import 'dart:developer';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

/**
 * This buttom sheet offers input fields to add a new book
 */

class AddBookBottomSheet extends StatelessWidget {
  const AddBookBottomSheet({super.key});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text("Go ahead and add another book from your physical bookshelf!"),
          ),
          SizedBox(height: 20),
          FractionallySizedBox(
            widthFactor: 0.8,
            child: SearchBar(),
          ),
          SizedBox(height: 20),
          FilledButton.icon(
              onPressed: () {
                log('pressed book scan button');
              },
              icon: Icon(Icons.camera_alt),
              label: Text("Scan ISBN code"))
        ],

      ),
    );
  }

}