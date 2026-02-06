import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:frontend/services/googlebooksapi_service.dart';
import 'package:frontend/services/library_service.dart';
import 'package:frontend/widgets/book_card.dart';
import 'package:provider/provider.dart';

import '../models/book.dart';

/// This bottom sheet offers input fields to add a new book
class AddBookBottomSheet extends StatefulWidget {
  const AddBookBottomSheet({super.key});

  @override
  State<AddBookBottomSheet> createState() => _AddBookBottomSheetState();
}

class _AddBookBottomSheetState extends State<AddBookBottomSheet> {
  Book? result;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: SingleChildScrollView(
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
                  try {
                    final fetchedBook = await GoogleBooksAPIService()
                        .searchByISBN(isbn);
                    setState(() {
                      result = fetchedBook;
                    });
                  } catch (e) {
                    log('ERROR: $e');
                  }
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
            if (result != null)
              InkWell(
                onTap: () => {
                  context.read<LibraryService>().addBook(result!),
                  print('Added book ${result?.title}')
                },
                child: SizedBox(
                    width: 200,
                    height: 280,
                    child: BookCard(book: result!, colorIndex: 1),
              )
              ),
          ],
        ),
      ),
    );
  }
}
