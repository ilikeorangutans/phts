import { TestBed } from '@angular/core/testing';

import { BasePathService } from './base-path.service';

describe('PathService', () => {
  let service: BasePathService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BasePathService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
