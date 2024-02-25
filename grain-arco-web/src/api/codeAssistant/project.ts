import axios from 'axios';
import { Base } from '@/api/base';

export interface Project {
    id: number;
    projectName: string;
    projectPath: string;
    webProjectPath: string;
    description: string;
}

export type CodeAssistants = Project[];


export interface PolicyParams extends Partial<Project> {
    page: number;
    pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
    data: CodeAssistants;
}

export function GetCodeAssistantList(data: any) {
    const uri = Object.entries(data)
        .map(
            ([key, value]) =>
                `${encodeURIComponent(key)}=${encodeURIComponent(value as string)}`
        )
        .join('&');
    return axios.get<PolicyListRes>(`/api/v1/codeAssistant/projects/list?${uri}`).then((res) => {
        return res.data;
    });
}

export function AddCodeAssistant(data: any) {
    return axios.post<Base>(`/api/v1/codeAssistant/projects`, data).then((res) => {
        return res.data;
    });
}

export function UpdateCodeAssistant(data: any) {
    return axios.put<Base>(`/api/v1/codeAssistant/projects`, data).then((res) => {
        return res.data;
    });
}

export function DeleteCodeAssistant(data: number) {
    return axios.delete<Base>(`/api/v1/codeAssistant/projects?pid=${data}`).then((res) => {
        return res.data;
    });
}
