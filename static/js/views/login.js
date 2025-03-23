import { Request } from "../lib/Request.js";

const do_login = async e => {
    e.preventDefault();

    const User = document.querySelector('[name="username"]').value.trim();
    const Password = document.querySelector('[name="password"]').value.trim();
    const ErrorNode = document.getElementById('ErrorNode');
    ErrorNode.classList.add('d-none');
    if(User === '' || Password === '') {
        ErrorNode.classList.remove('d-none');
        ErrorNode.textContent = 'Invalid credentials';
        return;
    }

    const { message } = await Request({ route: '/login', data: { User, Password } });
    if(message !== '') {
        ErrorNode.textContent = message;
        ErrorNode.classList.remove('d-none');
        return;
    }

};

document.getElementById('SignIn').addEventListener('click', do_login);