#!/bin/bash
set -e

echo "🚀 ebook-lcp-server Demo Setup - Pride and Prejudice"
echo "==================================================="

BOOK_DIR="examples/pride-and-prejudice"
INPUT="${BOOK_DIR}/pride-and-prejudice.epub"
OUTPUT="${BOOK_DIR}/pride-and-prejudice-protected.epub"
CONTENT_ID="pride-and-prejudice-001"

echo "🔐 Encrypting book (Content ID: $CONTENT_ID)..."

# More reliable way with explicit flags
lcpencrypt \
  -input "$INPUT" \
  -contentid "$CONTENT_ID" \
  -output "$OUTPUT" \
  -filename "pride-and-prejudice-protected.epub"

echo "✅ Encryption successful!"
echo "Protected file created → $OUTPUT"
echo ""
echo "Next steps:"
echo "   1. Start the server:     docker compose up -d"
echo "   2. Open Admin UI:        http://localhost:5173  (or your frontend port)"
echo "   3. Go to Publications → Import the protected EPUB"
echo "   4. Create a user"
echo "   5. Generate a license for the book"
echo ""
echo "📖 Test with Thorium Reader, Readium, or any LCP-compatible reader."
