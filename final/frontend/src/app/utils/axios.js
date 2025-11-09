// final/frontend/src/lib/axios.js  (ajustá la ruta según tu proyecto)
import axios from "axios";

// Base URL tomada de la variable pública de Next
// (si no está, usa localhost para desarrollo local)
const RAW_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Normalizamos: sin barra final
const BASE_URL = RAW_BASE.replace(/\/+$/, "");

const api = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
  // withCredentials: true, // habilítalo si usás cookies
});

// ===== Helpers =====
const ok = (r) => r.data;
const okResults = (r) => r.data?.results;

// ===== Cursos =====
export function search(query) {
  return api
    .get("/courses/search", { params: { query } })
    .then(ok)
    .catch((error) => {
      console.error("Hubo un Error en la busqueda:", error);
      throw error;
    });
}

export function getCourses() {
  return api
    .get("/courses")
    .then(okResults)
    .catch((error) => {
      console.error("Hubo un Error en la carga de cursos:", error);
      throw error;
    });
}

export function getCourseById(courseId) {
  return api
    .get(`/courses/${courseId}`)
    .then(ok)
    .catch((error) => {
      console.error("Hubo un Error en la obtencion del curso:", error);
      throw error;
    });
}

export function createCourse(courseRequest) {
  return api
    .post("/courses/create", courseRequest)
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en la creacion:", error);
      throw error;
    });
}

export function deleteCourse(courseId) {
  return api
    .delete(`/courses/delete/${courseId}`)
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en la eliminacion del curso:", error);
      throw error;
    });
}

export function updateCourse(courseId, updateRequest) {
  return api
    .put(`/courses/update/${courseId}`, updateRequest)
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en la actualizacion:", error);
      throw error;
    });
}

// ===== Usuarios / Auth =====
export function login(loginRequest) {
  return api
    .post("/users/login", loginRequest)
    .then((loginResponse) => {
      const token = loginResponse.data?.token;
      console.log("Token: ", token);
      localStorage.setItem("tokenType", token);
      localStorage.setItem("tokenId", token);
      return token;
    })
    .catch((error) => {
      console.log("Hubo un Error en el logueo:", error);
      throw error;
    });
}

export function registration(registrationRequest) {
  return api
    .post("/users/register", registrationRequest)
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en el registro:", error);
      throw error;
    });
}

export function getUserId(token) {
  return api
    .get("/users/userId", { headers: { Authorization: `Bearer ${token}` } })
    .then((r) => r.data?.message)
    .catch((error) => {
      console.error("Hubo un Error en la obtencion del id:", error);
      throw error;
    });
}

export function userAuthentication(token) {
  return api
    .get("/users/authentication", {
      headers: { Authorization: `Bearer ${token}` },
    })
    .then((r) => r.data?.message)
    .catch((error) => {
      console.error("Hubo un Error en la autenticación:", error);
      throw error;
    });
}

// ===== Subscripciones / Comentarios / Upload =====
export function subscribe(subscribeRequest) {
  return api
    .post("/subscriptions", subscribeRequest)
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en la subscripcion:", error);
      throw error;
    });
}

export function subscriptionList(userId) {
  return api
    .get(`/users/subscriptions/${userId}`)
    .then(okResults)
    .catch((error) => {
      console.error("Hubo un Error en la busqueda de inscripciones:", error);
      throw error;
    });
}

export function addComment(commentRequest) {
  return api
    .post("/users/comments", commentRequest)
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en en el comentado:", error);
      throw error;
    });
}

export function commentsList(courseId) {
  return api
    .get(`/courses/comments/${courseId}`)
    .then(okResults)
    .catch((error) => {
      console.error("Hubo un Error en la obtencion de los comentarios:", error);
      throw error;
    });
}

export function uploadFile(file, userId, courseId) {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("user_id", userId);
  formData.append("course_id", courseId);

  return api
    .post("/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    })
    .then(ok)
    .catch((error) => {
      console.log("Hubo un Error en la subida del archivo:", error);
      throw error;
    });
}