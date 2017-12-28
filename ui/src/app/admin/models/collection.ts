import { RenditionConfiguration } from './rendition-configuration';

export class Collection {
    name: string;
    id: number;
    slug: string;
    photoCount: number;
    createdAt: Date;
    updatedAt: Date;
    renditionConfigurations: Array<RenditionConfiguration>;
}
