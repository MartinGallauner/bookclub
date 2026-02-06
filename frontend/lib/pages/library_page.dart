import 'package:flutter/material.dart';
import 'package:frontend/main.dart';
import 'package:provider/provider.dart';

import '../models/book.dart';
import '../widgets/book_card.dart';

class LibraryPage extends StatelessWidget {
  const LibraryPage({super.key});

  @override
  Widget build(BuildContext context) {
    final appState = context.watch<AppState>();
    final List<Book> books = appState.libraryService.books;

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
