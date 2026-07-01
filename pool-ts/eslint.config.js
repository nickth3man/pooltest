import js from '@eslint/js';
import { defineConfig, globalIgnores } from 'eslint/config';
import tseslint from 'typescript-eslint';
import eslintConfigPrettier from 'eslint-config-prettier/flat';
import globals from 'globals';
import { importX } from 'eslint-plugin-import-x';

export default defineConfig([
  globalIgnores(['dist', 'node_modules']),

  js.configs.recommended,

  ...tseslint.configs.recommended,

  {
    files: ['src/**/*.ts'],
    languageOptions: {
      globals: {
        ...globals.browser,
      },
    },
    rules: {
      '@typescript-eslint/no-unused-vars': [
        'warn',
        { argsIgnorePattern: '^_' },
      ],
    },
  },

  {
    files: ['src/**/*.ts', 'vite.config.ts'],
    ...importX.flatConfigs.recommended,
    ...importX.flatConfigs.typescript,
    settings: {
      'import-x/resolver': {
        'eslint-import-resolver-typescript': {
          alwaysTryTypes: true,
          project: './tsconfig.json',
        },
      },
    },
  },

  eslintConfigPrettier,
]);
