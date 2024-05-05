import {Message} from "@arco-design/web-vue";

export const copyTextToClipboard = async (text: string) => {
    try {
        await navigator.clipboard.writeText(text);
        Message.success({
            content: '复制成功',
            duration: 5 * 1000,
        });
    } catch (err) {
        console.error('无法复制文本: ', err);
    }
}