"use strict";

const fs = require("fs");
const path = require("path");
const readFile = require("util").promisify(fs.readFile);

/** @type {import('sequelize-cli').Migration} */
module.exports = {
  async up(queryInterface) {
    try {
      const queryPath = path.resolve(__dirname, "../queries/create-tables.sql");
      const query = await readFile(queryPath, "utf8");
      return await queryInterface.sequelize.query(query);
    } catch (error) {
      console.error("Could not create tables", error)
    }
  },

  async down(queryInterface) {
      try {
      const queryPath = path.resolve(__dirname, "../queries/drop-tables.sql");
      const query = readFile(queryPath, "utf8");
      return await queryInterface.sequelize.query(query);
    } catch (error) {
      console.error("Could not create tables", error)
    }
  },
};
