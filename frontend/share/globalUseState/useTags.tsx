import { Tag } from "@/lib/api";
import { useState } from "react";

export const useTags = () => {

    const [allTags, setAllTags] = useState<Tag[]>([]);
    const [tagList, setTagList] = useState<Tag[]>([]);

    return {
        allTags,
        tagList,
        setAllTags,
        setTagList,
    }
}