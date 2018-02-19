import { Collection } from './collection';
export class Album {
  id: number;
  name: string;
  slug: string;
  photoCount: number;
  createdAt: Date;
  updatedAt: Date;
  collection: Collection;
}
