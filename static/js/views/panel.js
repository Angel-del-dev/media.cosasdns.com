import { Request } from "../lib/Request.js";
import { check_valid_domain } from "../lib/Token.js";
import { useState } from "../lib/hooks.js";

const [ getCurrentDirectory, setCurrentDirectory ] = useState('/');

const get_applications = async () => {
    const article = document.querySelector('article');
    article.innerHTML = '';
    const { message, applications } = await Request({ route: '/get-user-applications', data: {} });
    if(message !== '') throw new Error(message);
    
    applications.forEach((application, _) => {
        const node = document.createElement('div');
        node.classList.add('node', 'pointer');     
        node.title = `Application: ${application}`;   
        node.append(document.createTextNode(application));
        article.append(node);
    });
};

const create_application = async () => {
    const { message } = await Request({ route: '/create-application', data: {} });
    if(message !== '') throw new Error(message);
    get_applications();
};

const main = async () => {
    const ValidToken = await check_valid_domain();
    if(!ValidToken) return;    
    get_applications();
};

main();
document.getElementById('logout')?.addEventListener('click', _ => location.href = '/login');
document.getElementById('create-application')?.addEventListener('click', _ => create_application());