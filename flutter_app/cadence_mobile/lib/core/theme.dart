import 'package:flutter/material.dart';
class AppColors {
  // obsidian black, matte charcoal and Canary yellow pallete
  static const Color obsidianBlack = Color(0xFF0B0C0E);
  static const Color matteCharcoal = Color(0xFF1E2025);
  static const Color charcoalBorder = Color(0xFF2C2F36);

  static const Color warmCanaryYellow = Color(0xFFFFC700);
  static const Color warmCanaryYellowLight = Color(0xFFFFD84D);

  static const Color textPrimary = Color(0xFFF5F5F7);
  static const Color textSecondary = Color(0xFF8E95A3);
}
class AppTheme {
  static ThemeData get darkTheme{
    return ThemeData.dark().copyWith(
      scaffoldBackgroundColor: AppColors.matteCharcoal,
      primaryColor: AppColors.warmCanaryYellow,
      cardColor: AppColors.matteCharcoal,
      
      colorScheme: const ColorScheme.dark(
        primary: AppColors.obsidianBlack,
        secondary: AppColors.warmCanaryYellow,
        surface: AppColors.matteCharcoal,
        onPrimary: AppColors.obsidianBlack,
        onSurface: AppColors.textPrimary,
      ),

      appBarTheme: const AppBarTheme(
        backgroundColor: AppColors.obsidianBlack,
        elevation: 2,
        centerTitle: false,
        titleTextStyle: TextStyle(
          color: AppColors.textPrimary,
          fontSize: 20,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}