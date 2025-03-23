import { Request } from "../lib/Request.js";

export const check_valid_domain = async () => {
    const Token = localStorage.getItem('auth') ?? '';
    if(Token === '') return false;

    const { message } = await Request({ route: '/check-token', data: { Token } });
    
    if(message !== '') { location.href = '/login'; }
    return message === '';
};