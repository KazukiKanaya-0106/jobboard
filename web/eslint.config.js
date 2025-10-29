import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import pluginImport from "eslint-plugin-import";
import pluginUnusedImports from "eslint-plugin-unused-imports";
import prettier from "eslint-config-prettier";

export default tseslint.config(
  { files: ["**/*.{ts,tsx}"], ignores: ["dist/**", "node_modules/**"] },

  eslint.configs.recommended,
  tseslint.configs.strict,
  tseslint.configs.stylistic,

  {
    plugins: {
      import: pluginImport,
      "unused-imports": pluginUnusedImports,
    },
    rules: {
      "no-console": "off",
      "no-undef": "off",
      "import/extensions": "off",
      "unused-imports/no-unused-imports": "error",
      "unused-imports/no-unused-vars": [
        "warn",
        { vars: "all", varsIgnorePattern: "^_", args: "after-used", argsIgnorePattern: "^_" },
      ],
      "@typescript-eslint/no-unused-vars": ["error", { argsIgnorePattern: "^_", varsIgnorePattern: "^_" }],
      "@typescript-eslint/no-explicit-any": "error",
      "@typescript-eslint/triple-slash-reference": "off",
      "@typescript-eslint/no-unused-expressions": "off",
      "@typescript-eslint/consistent-type-definitions": ["error", "type"],
    },
  },

  prettier,
);