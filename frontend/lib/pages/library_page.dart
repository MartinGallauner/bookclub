import 'package:flutter/material.dart';
import 'package:frontend/services/library_service.dart';
import 'package:provider/provider.dart';

import '../models/book.dart';
import '../widgets/book_card.dart';

class LibraryPage extends StatelessWidget {
  const LibraryPage({super.key});

  @override
  Widget build(BuildContext context) {
    final libraryState = context.watch<LibraryService>();
    final List<Book> books = libraryState.books;

    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: GridView.builder(
        itemCount: books.length,
        gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
          maxCrossAxisExtent: 200,
          childAspectRatio: 0.7,
          crossAxisSpacing: 10,
          mainAxisSpacing: 10,
        ),
        itemBuilder: (context, index) {
          return BookCard(book: books[index], colorIndex: index);
        },
      ),
    );
  }
}
