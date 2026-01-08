import 'package:flutter/material.dart';

import '../models/book.dart';

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


class BookCard extends StatelessWidget {
  final Book book;
  final int colorIndex;
  const BookCard({
    super.key,
    required this.book,
    required this.colorIndex,
  });


  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Cover image
          Expanded(
            flex: 3,
            child: Container(
              width: double.infinity,
              color: Colors.primaries[colorIndex % Colors.primaries.length]
            )
          ),
          // Book details
          Padding(
              padding: EdgeInsets.all(8.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                book.title,
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              SizedBox(height: 4),
              Text(
                book.authors[0],
                style: TextStyle(
                  color: Colors.grey[600],
                  fontSize: 12,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              )
            ],
          ))
        ],
      ),

    );

  }

}