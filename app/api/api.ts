export function apiUrl(path:string){
    return "http://localhost:8080/v1"+(path.startsWith("/")?"":"/")+path
}