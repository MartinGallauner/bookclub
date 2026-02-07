class Book {
  final String isbn;
  final String title;
  final List<String> authors;
  final String coverUrl;
  final String? publisher;
  final int? publishedIn;
  final List<String>? genre;
  final String language;
  final String? description;

  const Book({
    required this.title,
    required this.authors,
    required this.isbn,
    required this.coverUrl,
    required this.language,
    this.publisher,
    this.publishedIn,
    this.genre,
    this.description,
  });
}
