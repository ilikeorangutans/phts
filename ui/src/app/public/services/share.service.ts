import { HttpClient } from '@angular/common/http';
import { Share } from './../models/share';
import { PathService } from './path.service';
import { Injectable } from '@angular/core';

@Injectable()
export class ShareService {

  constructor(
    private http: HttpClient,
    private pathService: PathService
  ) { }

  forSlug(slug: string): Promise<Share> {
    const url = this.pathService.shareBySlug(slug);
    return this.http.get<Share>(url).toPromise();
  }
}
