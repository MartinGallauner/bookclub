import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:frontend/services/googlebooksapi_service.dart';

/// This bottom sheet offers input fields to add a new book
class AddBookBottomSheet extends StatefulWidget {
  const AddBookBottomSheet({super.key});

  @override
  State<AddBookBottomSheet> createState() => _AddBookBottomSheetState();
}

class _AddBookBottomSheetState extends State<AddBookBottomSheet> {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              "Go ahead and add a book from your physical bookshelf!",
            ),
          ),
          SizedBox(height: 20),
          FractionallySizedBox(
            widthFactor: 0.8,
            child: TextField(
              onSubmitted: (isbn) async {
                final book = await GoogleBooksAPIService().searchByISBN(isbn);
                log(book.title);
              },
              decoration: InputDecoration(
                prefixIcon: Icon(Icons.search),
                hintText: 'Enter ISBN',
              ),
            ),
          ),
          SizedBox(height: 20),
          FilledButton.icon(
            onPressed: () {
              log('pressed book scan button');
            },
            icon: Icon(Icons.camera_alt),
            label: Text("Scan ISBN code"),
          ),
        ],
      ),
    );
  }
}
