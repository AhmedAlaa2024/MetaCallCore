use crate::{ c_void, c_char, CStr, c_int };

#[no_mangle]
pub extern "C" fn rs_loader_impl_load_from_file(_loader_impl: *mut c_void, paths: *const *mut c_char, size: usize) -> *mut c_void {
	// TODO: Not working
	/*
	for i in 0 .. size {
		let path: *const c_char = unsafe { *(paths.offset(i as isize)) };
		let c_path: &CStr = unsafe { CStr::from_ptr(path) };
		let path_slice: &str = c_path.to_str().unwrap();
		println!("RUSTY: {}", path_slice);
	}
	*/

	1 as c_int as *mut c_void
}
