import { defineConfig } from 'orval'

export default defineConfig({
  aiot: {
    input: {
      target: '../aiot-backend/docs/swagger.json',
    },
    output: {
      target: 'src/api/generated/index.ts',
      schemas: 'src/api/generated/model',
      client: 'react-query',
      httpClient: 'axios',
      override: {
        mutator: {
          path: './src/api/orval-mutator.ts',
          name: 'orvalAxios',
        },
      },
    },
  },
})
