import 'package:flutter/material.dart';

import '../models/book.dart';
import '../widgets/book_card.dart';

class LibraryPage extends StatelessWidget {
  const LibraryPage({super.key});

  @override
  Widget build(BuildContext context) {
    // Mock data for testing
    final List<Book> books = [
      Book(title: 'The Great Gatsby', authors: ['F. Scott Fitzgerald'], coverUrl: '', isbn: ''),
      Book(title: '1984', authors: ['George Orwell'], coverUrl: '', isbn: ''),
      Book(title: 'To Kill a Mockingbird', authors: ['Harper Lee'], coverUrl: '', isbn: ''),
      Book(title: 'Pride and Prejudice', authors: ['Jane Austen'], coverUrl: '', isbn: ''),
      Book(title: 'The Hobbit', authors: ['J.R.R. Tolkien'], coverUrl: '', isbn: ''),
      Book(title: 'Harry Potter', authors: ['J.K. Rowling'], coverUrl: '', isbn: ''),
    ];

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