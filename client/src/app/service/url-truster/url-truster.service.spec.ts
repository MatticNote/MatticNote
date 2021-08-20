import { TestBed } from '@angular/core/testing';

import { URLTrusterService } from './u-r-l-truster.service';

describe('UrlTrusterService', () => {
  let service: URLTrusterService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(URLTrusterService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
