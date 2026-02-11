class Profile {
  final String uid;
  final String email;
  final String displayName;
  final String photoURL;
  final DateTime updatedAt;
  final DateTime createdAt;

  const Profile({
    required this.uid,
    required this.email,
    required this.displayName,
    required this.photoURL,
    required this.updatedAt,
    required this.createdAt,
  });

  Map<String, dynamic> toMap() {
    return {
      'uid': uid,
      'email': email,
      'displayName': displayName,
      'photoURL': photoURL,
    };
  }
}
