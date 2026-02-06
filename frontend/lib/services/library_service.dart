import '../models/book.dart';

class LibraryService {
  final List<Book> _collection = [];

  LibraryService();

  List<Book> get books =>  _collection;

  void addBook(Book book) {
    _collection.add(book);
  }
}
