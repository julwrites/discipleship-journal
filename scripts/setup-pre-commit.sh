#!/bin/bash
# Setup pre-commit hooks for secret scanning

set -e

echo "🔧 Setting up pre-commit hooks for secret scanning..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Installing pre-commit..."
    pip install pre-commit || {
        echo "❌ Failed to install pre-commit. Please install it manually:"
        echo "   pip install pre-commit"
        echo "   or visit: https://pre-commit.com/#install"
        exit 1
    }
fi

# Install the git hooks
echo "📝 Installing git hooks..."
pre-commit install
pre-commit install --hook-type commit-msg
pre-commit install --hook-type pre-push

# Run initial scan
echo "🔍 Running initial secret scan..."
pre-commit run --all-files || {
    echo "⚠️  Initial scan found issues. Please review and fix them."
    echo "   To skip this check temporarily: git commit --no-verify"
    echo "   But please fix the issues before pushing to remote!"
}

echo "✅ Pre-commit hooks installed successfully!"
echo ""
echo "📋 Next steps:"
echo "   1. Review any issues found by the initial scan"
echo "   2. Commit your changes as usual"
echo "   3. The hooks will run automatically on each commit"
echo ""
echo "💡 Tips:"
echo "   - To skip hooks: git commit --no-verify"
echo "   - To run manually: pre-commit run --all-files"
echo "   - To update hooks: pre-commit autoupdate"