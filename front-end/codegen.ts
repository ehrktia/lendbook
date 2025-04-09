
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: "../internal/graph/schema.graphqls",
  generates: {
    "src/gql/": {
      preset: "client",
      plugins: ['typescript'],
      presetConfig:{
        gqlTagName: 'graphql',
      }
    }
  }
};

export default config;
