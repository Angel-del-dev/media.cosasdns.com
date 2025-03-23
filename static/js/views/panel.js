import { check_valid_domain } from "../lib/Token.js";

const main = async () => {
    await check_valid_domain();
    
    // console.log('Valid token');
};

main();