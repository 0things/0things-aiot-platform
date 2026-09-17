ARG LOGTO_TAG=latest
FROM svhd/logto:${LOGTO_TAG}

ENTRYPOINT ["sh", "-c", "npm run cli db seed -- --swe && npm start"]
